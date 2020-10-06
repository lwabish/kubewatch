package permission

import (
	"fmt"
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
	"github.com/bitnami-labs/kubewatch/pkg/utils"
	batchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Permission struct {
	ScName string
	chmod string
	chown string
}

func (p *Permission) Init(c *config.Config) error {
	p.ScName=c.Handler.Permission.ScName
	p.chmod=c.Handler.Permission.Chmod
	p.chown=c.Handler.Permission.Chown
	//TODO 环境变量
	//TODO 空值检查
	return nil
}

func (p *Permission) Handle(e event.Event)  {
	if e.Reason!="Created" || e.Kind!="persistent volume" {
		//fmt.Println("跳过:",e.Reason,e.Kind)
		return
	}
	//fmt.Println(p.ScName,p.chmod,p.chown)

	var clientSet kubernetes.Interface

	if _, err := rest.InClusterConfig(); err != nil {
		clientSet = utils.GetClientOutOfCluster()
	} else {
		clientSet = utils.GetClient()
	}

	pv, _ := clientSet.CoreV1().PersistentVolumes().Get(e.Name, metaV1.GetOptions{})
	fmt.Println(pv.Spec.ClaimRef.Name, pv.Spec.ClaimRef.Namespace, pv.Status.Phase)

	newJob := &batchV1.Job{
		ObjectMeta: metaV1.ObjectMeta{
			Name: fmt.Sprintf("pfixer-%s-%s",pv.Spec.ClaimRef.Name,pv.Spec.ClaimRef.Namespace),
			Namespace: pv.Spec.ClaimRef.Namespace,
			Labels: map[string]string{"jobgroup":"permission"},
		},
		Spec: batchV1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Name: fmt.Sprintf("pfixer-%s-%s",pv.Spec.ClaimRef.Name,pv.Spec.ClaimRef.Namespace),
					Labels: map[string]string{"jobgroup":"permission"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "permission-fixer",
							Image: "alpine",
							Command: []string{
								"sh",
								"-c",
								fmt.Sprintf("chmod -R %s /data && chown -R %s:%s /data",p.chmod,p.chown,p.chown),
							},
							VolumeMounts: []v1.VolumeMount{
								{
									MountPath: "/data",
									Name: "data",
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
					Volumes: []v1.Volume{
						{
							Name: "data",
							VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: pv.Spec.ClaimRef.Name,
								},
							},
						},
					},
				},
			},
		},
	}
	result,_:=clientSet.BatchV1().Jobs(pv.Spec.ClaimRef.Namespace).Create(newJob)
	fmt.Println(result)
}