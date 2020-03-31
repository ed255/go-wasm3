# Android build

```
set -ex

mkdir -p /tmp/gomobile
docker run -it \
    --mount type=bind,source=$HOME/git,target=/root/git,readonly \
    --mount type=bind,source=/tmp/gomobile,target=/tmp/gomobile \
    android-build-box \
    /bin/bash -c 'set -ex && \
        apt-get update && \
        apt-get install -y cmake && \
        rm -r /tmp/gomobile/git || true && \
        mkdir -p /tmp/gomobile/git && \
        cp -r /root/git/go-wasm3 /tmp/gomobile/git && \
        cp -r /root/git/wasm3 /tmp/gomobile/git && \
        cd /tmp/gomobile/git/wasm3 && \
        cd android && \
        rm -r build || true && \
        mkdir build && \
        cd build && \
        NDK_HOME=/opt/android-ndk/android-ndk-r20 ../make_all.sh && \
        cp -r android /tmp/gomobile/git/go-wasm3/lib/
        cd /tmp/gomobile/git/go-wasm3 && \
        go build && \
        gomobile bind -androidapi=29 --target android -o /tmp/gomobile/wasm.aar'
    # /bin/bash
```
