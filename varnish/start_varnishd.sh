#!/bin/sh
echo $VARNISH_SECRET > $VCL_DIR/secret

# Start varnish and log
# echo "** Varnish cache size is set to ${VARNISH_MEMORY}"
# echo "** Starting Varnish cache with no backends"

varnishd -F -s malloc,${VARNISH_MEMORY} \
          -a 0.0.0.0:${VARNISH_PORT} \
          -p default_ttl=3600 \
          -p default_grace=3600 \
          -S $VCL_DIR/secret \
          -T 127.0.0.1:6082 \
          -f $VCL_DIR/$DEFAULT_FILE
