# How to Install file2byteslice:
# go get github.com/hajimehoshi/file2byteslice
# go install github.com/hajimehoshi/file2byteslice

INPUT_FILE=
OUTPUT_FILE=
PACKAGE_NAME=
VARIABLE_NAME=

generate () {
  DIRECTORY=$1
  PACKAGE_NAME=$2
  for f in `ls $DIRECTORY`; do
    echo "processing $f"
    INPUT_FILE=$DIRECTORY/$f
    OUTPUT_FILE=$DIRECTORY/${f%.*}.go
    VARIABLE_NAME=`echo ${f%.*} | tr '[:lower:]' '[:upper:]'`
    ~/go/bin/file2byteslice -input $INPUT_FILE -output $OUTPUT_FILE -package $PACKAGE_NAME -var $VARIABLE_NAME
  done
}

generate "gfx" "gfx"

