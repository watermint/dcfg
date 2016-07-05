
export TARGET_OS="windows darwin linux"
BUILD=`pwd`/out
DIST=/dist
APP_VERSION=`cat version`
echo building: $APP_VERSION
echo UID: `id`

for t in $TARGET_OS; do
  mkdir -p "$BUILD/$t";
done

for t in $TARGET_OS; do
  echo Building: $t
  GOOS=$t GOARCH=amd64 go build -ldflags "-X main.AppVersion=`cat version`" -o "$BUILD/$t/dcfg" github.com/watermint/dcfg
done

mv $BUILD/windows/dcfg $BUILD/windows/dcfg.exe # workaround
cd $BUILD
zip -9 -r $DIST/dcfg-$APP_VERSION.zip .
