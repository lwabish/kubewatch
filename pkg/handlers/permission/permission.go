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
	"log"
	"os"
)

type Permission struct {
	scName string
	chmod  string
	chown  string
}

func (p *Permission) Init(c *config.Config) error {
	p.scName = os.Getenv("scName")
	p.chown = os.Getenv("chown")
	p.chmod = os.Getenv("chmod")

	if p.scName == "" {
		p.scName = c.Handler.Permission.ScName
	}
	if p.chmod == "" {
		p.chmod = c.Handler.Permission.Chmod
		if p.chmod == "" {
			p.chmod = "777"
		}
	}
	if p.chown == "" {
		p.chown = c.Handler.Permission.Chown
		if p.chown == "" {
			p.chown = "root"
		}
	}

	return nil
}

func (p *Permission) Handle(e event.Event) {
	if e.Reason != "Created" || e.Kind != "persistent volume" {
		return
	}

	var clientSet kubernetes.Interface

	if _, err := rest.InClusterConfig(); err != nil {
		clientSet = utils.GetClientOutOfCluster()
	} else {
		clientSet = utils.GetClient()
	}

	pv, err := clientSet.CoreV1().PersistentVolumes().Get(e.Name, metaV1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if pv.Spec.StorageClassName != p.scName {
		return
	}
	//log.Println(pv.Spec.ClaimRef.Name, pv.Spec.ClaimRef.Namespace, pv.Status.Phase)

	newJob := &batchV1.Job{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      fmt.Sprintf("pfixer-%s-%s", pv.Spec.ClaimRef.Name, pv.Spec.ClaimRef.Namespace),
			Namespace: pv.Spec.ClaimRef.Namespace,
			Labels:    map[string]string{"jobgroup": "permission"},
		},
		Spec: batchV1.JobSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Name:   fmt.Sprintf("pfixer-%s-%s", pv.Spec.ClaimRef.Name, pv.Spec.ClaimRef.Namespace),
					Labels: map[string]string{"jobgroup": "permission"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "permission-fixer",
							Image: "alpine",
							Command: []string{
								"sh",
								"-c",
								fmt.Sprintf("chmod -R %s /data && chown -R %s:%s /data", p.chmod, p.chown, p.chown),
							},
							VolumeMounts: []v1.VolumeMount{
								{
									MountPath: "/data",
									Name:      "data",
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
	_, err = clientSet.BatchV1().Jobs(pv.Spec.ClaimRef.Namespace).Create(newJob)
	if err != nil {
		log.Fatal(err)
	}
}
