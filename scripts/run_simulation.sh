SCRIPTPATH="$(
  cd "$(dirname "$0")"
  pwd -P
)"

echo "Running simulation $1..."

go test -mod=readonly github.com/ixoworld/bonds/app \
  -run=TestFullAppSimulation \
  -v \
  -timeout 24h \
  -Enabled=true \
  -NumBlocks=100 \
  -BlockSize=200 \
  -Commit=true \
  -Seed="$1" \
  -Period=5 \
  -Verbose \
  -Params="$SCRIPTPATH"/../input/SimulationParamsFile.json \
  -ExportStatePath="$SCRIPTPATH"/../ExportState"$1$2".json \
  -ExportStatsPath="$SCRIPTPATH"/../ExportStats"$1$2".json \
  -PrintAllInvariants >"$SCRIPTPATH"/../ExportLog"$1$2".txt

echo "...$1 DONE."