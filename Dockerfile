FROM alpine:3.5

ENV SOURCE_DIR /image-resolver-src

COPY . $SOURCE_DIR/

RUN apk add --no-cache  --update bash ca-certificates \
  && apk --no-cache --virtual .build-dependencies add git go libc-dev \
  && cd $SOURCE_DIR \
  && BUILDINFO_PACKAGE="github.com/Financial-Times/service-status-go/buildinfo." \
  && VERSION="version=$(git describe --tag --always 2> /dev/null)" \
  && DATETIME="dateTime=$(date -u +%Y%m%d%H%M%S)" \
  && REPOSITORY="repository=$(git config --get remote.origin.url)" \
  && REVISION="revision=$(git rev-parse HEAD)" \
  && BUILDER="builder=$(go version)" \
  && LDFLAGS="-X '"${BUILDINFO_PACKAGE}$VERSION"' -X '"${BUILDINFO_PACKAGE}$DATETIME"' -X '"${BUILDINFO_PACKAGE}$REPOSITORY"' -X '"${BUILDINFO_PACKAGE}$REVISION"' -X '"${BUILDINFO_PACKAGE}$BUILDER"'" \
  && cd .. \
  && export GOPATH=/gopath \
  && go get -u github.com/kardianos/govendor \
  && REPO_PATH="github.com/Financial-Times/image-resolver" \
  && mkdir -p $GOPATH/src/${REPO_PATH} \
  && cp -r $SOURCE_DIR/* $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && $GOPATH/bin/govendor sync \
  && go get ./... \
  && echo ${LDFLAGS} \
  && go build -ldflags="${LDFLAGS}" \
  && mv image-resolver / \
  && apk del .build-dependencies \
  && rm -rf $GOPATH /var/cache/apk/*

CMD [ "/image-resolver" ]
