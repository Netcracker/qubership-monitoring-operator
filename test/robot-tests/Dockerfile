# Based on python:3.13.2-alpine3.21
ARG BUILDPLATFORM
FROM --platform=$BUILDPLATFORM ghcr.io/netcracker/qubership-docker-integration-tests:0.4.3

# ROBOT_HOME is defined in the base image.
ENV ROBOT_OUTPUT=${ROBOT_HOME}/output \
    ROBOT_UID=1000 \
    ROBOT_GID=0 \
    SERVICE_CHECKER_SCRIPT=${ROBOT_HOME}/monitoring_pods_checker.py \
    SERVICE_CHECKER_SCRIPT_TIMEOUT=500

COPY --chown=${ROBOT_UID}:${ROBOT_GID} requirements.txt ${ROBOT_HOME}/requirements.txt
COPY --chown=${ROBOT_UID}:${ROBOT_GID} src/ ${ROBOT_HOME}/

RUN set -x \
    && pip3 install --no-cache-dir -r ${ROBOT_HOME}/requirements.txt \
    && rm -rf /var/cache/apk/*

USER ${ROBOT_UID}:${ROBOT_GID}

EXPOSE 8080
VOLUME ["${ROBOT_OUTPUT}"]
