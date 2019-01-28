#!/bin/sh

if [ -z "${COMPUTE_RANK_ON_GPU}" ]
then
  COMPUTE_GPU=true
else
  COMPUTE_GPU="${COMPUTE_RANK_ON_GPU}"
fi

echo ${COMPUTE_GPU}
cyberd start --compute-rank-on-gpu=${COMPUTE_GPU}
