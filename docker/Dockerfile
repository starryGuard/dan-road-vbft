
FROM tianon/ubuntu-core:14.04
ENV DAN_PATH /var/dan-road-vbft
RUN mkdir -p $DAN_PATH
COPY dan-road-vbft $DAN_PATH
RUN chmod +x /var/dan-road-vbft/dan-road-vbft
EXPOSE 20334 20335 20336 20337 20338 20339
WORKDIR $DAN_PATH
ENTRYPOINT ["./dan-road-vbft"]
