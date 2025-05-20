echo $GOOS $GOARCH
if [[ "$GOOS" == "linux" ]]; then
  echo "Deleting rsrc.syso for linux build"
  rm -f rsrc.syso
fi