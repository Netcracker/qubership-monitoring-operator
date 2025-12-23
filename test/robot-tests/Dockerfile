# Based on python:3.13.2-alpine3.21
ARG BUILDPLATFORM
FROM --platform=$BUILDPLATFORM ghcr.io/netcracker/qubership-docker-integration-tests:0.3.5

ENV ROBOT_OUTPUT=/opt/robot/output \
    ROBOT_UID=1000 \
    ROBOT_GID=0 \
    SERVICE_CHECKER_SCRIPT=${ROBOT_HOME}/monitoring_pods_checker.py \
    SERVICE_CHECKER_SCRIPT_TIMEOUT=500

# Copy tests parts
COPY requirements.txt ${ROBOT_HOME}/requirements.txt
COPY src/ ${ROBOT_HOME}/

RUN \
    # Install python requirements
    pip3 install --no-cache-dir -r ${ROBOT_HOME}/requirements.txt \
    # Create and grant permissions for tests output
    && mkdir -p ${ROBOT_OUTPUT} \
    && chmod a+w ${ROBOT_OUTPUT}

USER ${ROBOT_UID}:${ROBOT_GID}

EXPOSE 8080
VOLUME ["${ROBOT_OUTPUT}"]
