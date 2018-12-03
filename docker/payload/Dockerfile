
#FROM tianon/ubuntu-core:14.04
#ENV DAN_PATH /var/dan-road-vbft
#RUN mkdir -p $DAN_PATH
#COPY / $DAN_PATH
#RUN chmod +x /var/dan-road-vbft/run.sh
#EXPOSE 20334 20335 20336 20337 20338 20339
#WORKDIR $DAN_PATH
#ENTRYPOINT ["./run.sh"]


FROM tianon/ubuntu-core:14.04
ENV CVBFT_PATH /var/cvbft
RUN mkdir -p $CVBFT_PATH
COPY cvbft $CVBFT_PATH
COPY wallet.dat $CVBFT_PATH
EXPOSE 20334 20335 20336 20337 20338 20339
WORKDIR $CVBFT_PATH
ENTRYPOINT ["./cvbft"]


