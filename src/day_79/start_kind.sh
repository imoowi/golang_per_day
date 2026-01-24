#创建kind集群
kind create cluster --config kind-config.yaml --name golang-per-day
#检查节点是否Ready
kubectl get nodes
# Setting LINKERD2_VERSION sets the version to install.
# If unset, you'll get the latest available edge version.
export LINKERD2_VERSION=edge-25.10.7
curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install-edge | sh

export PATH=$HOME/.linkerd2/bin:$PATH

linkerd version

kubectl apply --server-side -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.4.0/standard-install.yaml 
linkerd check --pre

linkerd install --crds | kubectl apply -f -
linkerd install | kubectl apply -f -

kubectl -n linkerd get deploy
linkerd check
kubectl apply -f https://run.linkerd.io/emojivoto.yml
kubectl -n emojivoto get deploy -o yaml \
  | linkerd inject - \
  | kubectl apply -f -

kubectl -n emojivoto get pods

linkerd viz install | kubectl apply -f -
linkerd viz dashboard
