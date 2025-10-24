cp /mnt/c/Users/Rahul\ Poonia/Downloads/extractous-ffi-linux_amd64.zip ~/dev/extractous-go/

unzip extractous-ffi-linux_amd64.zip -d extractous-ffi-linux_amd64

cd extractous-ffi-linux_amd64

cp ./lib ../../benchmark/native/linux_amd64/lib -r
