
export TARGET_OS="windows darwin linux"
DIST=/dist
APP_VERSION=`cat version`
echo building: $APP_VERSION
echo UID: `id`

for t in $TARGET_OS; do
  mkdir -p "$DIST/$t";
done

for t in $TARGET_OS; do
  echo Building: $t
  GOOS=$t GOARCH=amd64 go build -ldflags "-X main.AppVersion=`cat version`" -o "$DIST/$t/dcfg" github.com/watermint/dcfg
done

mv $DIST/windows/dcfg $DIST/windows/dcfg.exe # workaround
cd $DIST
zip -9 -r dcfg-$APP_VERSION.zip .
