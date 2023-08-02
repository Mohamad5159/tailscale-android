# This is a Dockerfile for creating a build environment for
# tailscale-android.

FROM --platform=linux/amd64 eclipse-temurin:20-jdk

# To enable running android tools such as aapt
RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y libz1 libstdc++6 unzip
# For Go:
RUN apt-get -y --no-install-recommends install curl gcc
RUN apt-get -y --no-install-recommends install ca-certificates libc6-dev git

RUN apt-get -y install make

RUN mkdir -p build
ENV HOME /build

# Make android sdk location, the later make step will populate it.
RUN mkdir android-sdk
ENV ANDROID_HOME $HOME/android-sdk
ENV ANDROID_SDK_ROOT $ANDROID_HOME
ENV PATH $PATH:$HOME/bin:$ANDROID_HOME/platform-tools

# We need some version of Go new enough to support the "embed" package
# to run "go run tailscale.com/cmd/printdep" to figure out which Tailscale Go
# version we need later, but otherwise this toolchain isn't used:
RUN curl -L https://go.dev/dl/go1.20.5.linux-amd64.tar.gz | tar -C /usr/local -zxv
RUN ln -s /usr/local/go/bin/go /usr/bin

RUN mkdir -p $HOME/tailscale-android
RUN git config --global --add safe.directory $HOME/tailscale-android
WORKDIR $HOME/tailscale-android

# Preload Android SDK
COPY Makefile Makefile
# Get android sdk, ndk, and rest of the stuff needed to build the android app.
RUN make androidsdk

# Preload Gradle
COPY android/gradlew android/gradlew
COPY android/gradle android/gradle
RUN ./android/gradlew

CMD /bin/bash
