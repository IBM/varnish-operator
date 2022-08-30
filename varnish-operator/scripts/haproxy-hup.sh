#!/bin/bash
kill -HUP "$(pgrep -o haproxy)"
