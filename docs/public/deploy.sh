#!/bin/sh

# if has change, return, otherwise exit this script with 0
function is_doc_changed() {
    last_commit=$(ssh ${WEB_HOST_USER}@${WEB_HOST} "cat ${WEB_DEST_DIR}/last_commit")
    if [ $? -ne 0 ]; then
        echo "ERROR: Failed to read last commit, always generate and deploy doc."
        return
    fi

    git diff --name-only ${last_commit}    
    if [ $? -ne 0 ]; then
        echo "ERROR: Unknown commit hash, always generate and deploy doc."
        return
    fi

    git diff --name-only ${last_commit} | grep -e "^docs\/public\/"
    if [ $? -ne 0 ]; then
        echo "INFO: No new commit for public docs since commit ${last_commit}."
        exit 0
    fi

    # generate and deploy docs here
}


script_dir=$(dirname $0)
script_abs_dir=$(realpath ${script_dir})
commit_short_hash=$(git rev-parse --short HEAD)
commit_hash=$(git rev-parse HEAD)

is_doc_changed

# web pages were already built here, and can be found under speedle/site
# archive and compress it
cd ${script_abs_dir}/speedle/site
tar cvzf ../${commit_short_hash}.tar.gz *
cd ${script_abs_dir}/speedle

# ${WEB_HOST} is the host name, ${WEB_DEST_DIR} is the web page dir, ${WEB_HOST_USER} is the login name
# The two env vars were set by the runner

# generate script for sftp
cat << EOF > /tmp/sftp.script
cd ${WEB_DEST_DIR}
put ${commit_short_hash}.tar.gz
EOF

sftp -b /tmp/sftp.script ${WEB_HOST_USER}@${WEB_HOST}

ssh ${WEB_HOST_USER}@${WEB_HOST} "rm -rf ${WEB_DEST_DIR}/site/*; tar -C ${WEB_DEST_DIR}/site -xvf ${WEB_DEST_DIR}/${commit_short_hash}.tar.gz; echo ${commit_hash} > ${WEB_DEST_DIR}/last_commit"

