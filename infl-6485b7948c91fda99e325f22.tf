apiVersion: storage.k8s.io/v1
kind: VolumeAttachment
metadata:
  annotations:
    csi.alpha.kubernetes.io/node-id: i-002e6d0a53f79b89f
  name: csi-a9790b1f1ff63d48bd4d6a92e47c4e07927d9e03c28202a17f2ffb642921747a
spec:
  attacher: ebs.csi.aws.com
  source:
    persistentVolumeName: pvc-230eddae-6209-41cd-98e9-8d61d21ad5b5
