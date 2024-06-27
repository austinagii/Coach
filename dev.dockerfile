FROM docker:cli 

RUN apk update && apk upgrade && apk add bash

WORKDIR /coach

RUN ln -s /coach/coach.sh /usr/local/bin/coach

ENTRYPOINT ["bash"]
