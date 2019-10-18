package chaosengine

import (
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Determine whether apps with matching labels have chaos annotation set to true
func checkAnnotationDaemonSet(clientSet *kubernetes.Clientset) error {

	targetApplicationDaemonSetList, err := clientSet.AppsV1().DaemonSets(engine.appInfo.namespace).List(metav1.ListOptions{LabelSelector: engine.instance.Spec.Appinfo.Applabel, FieldSelector: ""})
	if err != nil {
		log.Error(err, "unable to list apps matching labels")
		return err
	}
	chaosCandidates := 0
	if len(targetApplicationDaemonSetList.Items) == 0 {
		return errors.New("No app daemonsets with matching labels")
	}
	for _, app := range targetApplicationDaemonSetList.Items {
		engine.appName = app.ObjectMeta.Name
		engine.appUUID = app.ObjectMeta.UID
		//Checks if the annotation is "true" / "false"
		annotationValue := app.ObjectMeta.GetAnnotations()[chaosAnnotationKey]

		if annotationValue == chaosAnnotationValue {
			// Add it to the Chaos Candidates, and log the details
			log.Info("chaos candidate : ", "appName", engine.appName, "appUUID", engine.appUUID)
			chaosCandidates++
		}
		if chaosCandidates > 1 {
			return errors.New("Too many chaos candidates with same label, either provide unique labels or annotate only desired app for chaos")
		}
		if chaosCandidates == 0 {
			return errors.New("no chaos-candidate found")

		}
	}

	return nil
}
