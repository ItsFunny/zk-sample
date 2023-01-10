FROM okexchain/build-env:go1.18.5-static as okexchain-builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct
WORKDIR /root
COPY exchain ./exchain
RUN cd exchain && make install WITH_ROCKSDB=true VenusHeight=1 LINK_STATICALLY=true && rm -rf /go/pkg && rm -rf .cache && cd /root && rm -rf /root/exchain

VOLUME ["/data/lrp_node"]
CMD ["exchaind"]
