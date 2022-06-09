package v1alpha1

// import (
// 	"testing"

// 	"gopkg.in/yaml.v2"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// )

// func setup(t *testing.T) {
// 	obj := make(map[string]interface{})
// 	err := yaml.Unmarshal([]byte(prYaml), &obj)
// 	if err != nil {
// 		t.Errorf("failed to unmarshal test pr: %s", err)
// 		t.FailNow()
// 	}
// 	uns = &unstructured.Unstructured{
// 		Object: obj,
// 	}
// }

// func TestTaskRunGraph(t *testing.T) {
// 	setup(t)
// 	pr, err := NewPipelineRun(uns)
// 	if err != nil {
// 		t.Errorf("failed to create new pr: %s", err)
// 		t.FailNow()
// 	}

// 	if pr.TotalTasks != 3 {
// 		t.Errorf("found %d total tasks; should be 3", pr.TotalTasks)
// 	}

// }

// var uns *unstructured.Unstructured

// var prYaml string = `apiVersion: tekton.dev/v1beta1
// kind: PipelineRun
// metadata:
//   creationTimestamp: "2022-06-09T17:20:54Z"
//   generateName: myapp-
//   generation: 1
//   labels:
//     tekton.dev/pipeline: myapp-4zhw4
//   name: myapp-4zhw4
//   namespace: tekton-pipelines
//   resourceVersion: "1358701"
//   uid: 98fd81f9-58cb-4107-a4ac-49a26a368f7e
// spec:
//   params:
//   - name: MESSAGE
//     value: Hi!
//   pipelineSpec:
//     params:
//     - name: MESSAGE
//       type: string
//     tasks:
//     - name: clone
//       params:
//       - name: url
//         value: https://github.com/lilley2412/samples-go-api.git
//       - name: revision
//         value: main
//       taskRef:
//         kind: Task
//         name: git-clone
//       workspaces:
//       - name: output
//         subPath: $(context.pipelineRun.uid)
//         workspace: build
//     - name: build
//       runAfter:
//       - clone
//       taskSpec:
//         metadata: {}
//         spec: null
//         steps:
//         - image: golang:1.18.2-alpine
//           name: pkg
//           resources: {}
//           script: |
//             #!/bin/sh
//             cd "$(workspaces.output.path)"
//             go mod download -x
//           volumeMounts:
//           - mountPath: /go/pkg
//             name: go-mod
//           - mountPath: /root/.cache/go-build
//             name: go-build
//         - image: golang:1.18.2-alpine
//           name: build
//           resources: {}
//           script: |
//             #!/bin/sh
//             cd "$(workspaces.output.path)"
//             CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main
//             echo "created go binary"
//           volumeMounts:
//           - mountPath: /go/pkg
//             name: go-mod
//           - mountPath: /root/.cache/go-build
//             name: go-build
//         volumes:
//         - name: go-mod
//           persistentVolumeClaim:
//             claimName: go-mod
//         - name: go-build
//           persistentVolumeClaim:
//             claimName: go-build
//         workspaces:
//         - name: output
//       workspaces:
//       - name: output
//         subPath: $(context.pipelineRun.uid)
//         workspace: build
//     - name: containerize
//       runAfter:
//       - build
//       taskSpec:
//         metadata: {}
//         spec: null
//         steps:
//         - image: gcr.io/kaniko-project/executor:debug
//           name: containerize
//           resources: {}
//           script: |
//             #!/busybox/sh

//             /kaniko/executor --context=dir://$(workspaces.output.path) \
//               --snapshotMode=redo \
//               --insecure-registry 10.0.10.1:30500 \
//               --dockerfile Dockerfile \
//               --destination 10.0.10.1:30500/samples-go-api:latest
//         workspaces:
//         - name: output
//       workspaces:
//       - name: output
//         subPath: $(context.pipelineRun.uid)
//         workspace: build
//     workspaces:
//     - name: build
//   serviceAccountName: default
//   timeout: 1h0m0s
//   workspaces:
//   - name: build
//     persistentVolumeClaim:
//       claimName: build
// status:
//   completionTime: "2022-06-09T17:21:12Z"
//   conditions:
//   - lastTransitionTime: "2022-06-09T17:21:12Z"
//     message: 'Tasks Completed: 3 (Failed: 0, Cancelled 0), Skipped: 0'
//     reason: Succeeded
//     status: "True"
//     type: Succeeded
//   pipelineSpec:
//     params:
//     - name: MESSAGE
//       type: string
//     tasks:
//     - name: clone
//       params:
//       - name: url
//         value: https://github.com/lilley2412/samples-go-api.git
//       - name: revision
//         value: main
//       taskRef:
//         kind: Task
//         name: git-clone
//       workspaces:
//       - name: output
//         subPath: $(context.pipelineRun.uid)
//         workspace: build
//     - name: build
//       runAfter:
//       - clone
//       taskSpec:
//         metadata: {}
//         spec: null
//         steps:
//         - image: golang:1.18.2-alpine
//           name: pkg
//           resources: {}
//           script: |
//             #!/bin/sh
//             cd "$(workspaces.output.path)"
//             go mod download -x
//           volumeMounts:
//           - mountPath: /go/pkg
//             name: go-mod
//           - mountPath: /root/.cache/go-build
//             name: go-build
//         - image: golang:1.18.2-alpine
//           name: build
//           resources: {}
//           script: |
//             #!/bin/sh
//             cd "$(workspaces.output.path)"
//             CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main
//             echo "created go binary"
//           volumeMounts:
//           - mountPath: /go/pkg
//             name: go-mod
//           - mountPath: /root/.cache/go-build
//             name: go-build
//         volumes:
//         - name: go-mod
//           persistentVolumeClaim:
//             claimName: go-mod
//         - name: go-build
//           persistentVolumeClaim:
//             claimName: go-build
//         workspaces:
//         - name: output
//       workspaces:
//       - name: output
//         subPath: $(context.pipelineRun.uid)
//         workspace: build
//     - name: containerize
//       runAfter:
//       - build
//       taskSpec:
//         metadata: {}
//         spec: null
//         steps:
//         - image: gcr.io/kaniko-project/executor:debug
//           name: containerize
//           resources: {}
//           script: |
//             #!/busybox/sh

//             /kaniko/executor --context=dir://$(workspaces.output.path) \
//               --snapshotMode=redo \
//               --insecure-registry 10.0.10.1:30500 \
//               --dockerfile Dockerfile \
//               --destination 10.0.10.1:30500/samples-go-api:latest
//         workspaces:
//         - name: output
//       workspaces:
//       - name: output
//         subPath: $(context.pipelineRun.uid)
//         workspace: build
//     workspaces:
//     - name: build
//   startTime: "2022-06-09T17:20:54Z"
//   taskRuns:
//     myapp-4zhw4-build:
//       pipelineTaskName: build
//       status:
//         completionTime: "2022-06-09T17:21:06Z"
//         conditions:
//         - lastTransitionTime: "2022-06-09T17:21:06Z"
//           message: All Steps have completed executing
//           reason: Succeeded
//           status: "True"
//           type: Succeeded
//         podName: myapp-4zhw4-build-pod
//         startTime: "2022-06-09T17:21:00Z"
//         steps:
//         - container: step-pkg
//           imageID: docker.io/library/golang@sha256:4795c5d21f01e0777707ada02408debe77fe31848be97cf9fa8a1462da78d949
//           name: pkg
//           terminated:
//             containerID: containerd://d56303dc9518ea95193897c0e9a56fb5d4bd5b6c35be0d5698bbc9b72ac98011
//             exitCode: 0
//             finishedAt: "2022-06-09T17:21:03Z"
//             reason: Completed
//             startedAt: "2022-06-09T17:21:03Z"
//         - container: step-build
//           imageID: docker.io/library/golang@sha256:4795c5d21f01e0777707ada02408debe77fe31848be97cf9fa8a1462da78d949
//           name: build
//           terminated:
//             containerID: containerd://bad520646c9a3cf55ce38f422813db9d18684212c6bf1e88ec22a507fde46b03
//             exitCode: 0
//             finishedAt: "2022-06-09T17:21:03Z"
//             reason: Completed
//             startedAt: "2022-06-09T17:21:03Z"
//         taskSpec:
//           steps:
//           - image: golang:1.18.2-alpine
//             name: pkg
//             resources: {}
//             script: |
//               #!/bin/sh
//               cd "$(workspaces.output.path)"
//               go mod download -x
//             volumeMounts:
//             - mountPath: /go/pkg
//               name: go-mod
//             - mountPath: /root/.cache/go-build
//               name: go-build
//           - image: golang:1.18.2-alpine
//             name: build
//             resources: {}
//             script: |
//               #!/bin/sh
//               cd "$(workspaces.output.path)"
//               CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main
//               echo "created go binary"
//             volumeMounts:
//             - mountPath: /go/pkg
//               name: go-mod
//             - mountPath: /root/.cache/go-build
//               name: go-build
//           volumes:
//           - name: go-mod
//             persistentVolumeClaim:
//               claimName: go-mod
//           - name: go-build
//             persistentVolumeClaim:
//               claimName: go-build
//           workspaces:
//           - name: output
//     myapp-4zhw4-clone:
//       pipelineTaskName: clone
//       status:
//         completionTime: "2022-06-09T17:21:00Z"
//         conditions:
//         - lastTransitionTime: "2022-06-09T17:21:00Z"
//           message: All Steps have completed executing
//           reason: Succeeded
//           status: "True"
//           type: Succeeded
//         podName: myapp-4zhw4-clone-pod
//         startTime: "2022-06-09T17:20:54Z"
//         steps:
//         - container: step-clone
//           imageID: gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/git-init@sha256:c0b0ed1cd81090ce8eecf60b936e9345089d9dfdb6ebdd2fd7b4a0341ef4f2b9
//           name: clone
//           terminated:
//             containerID: containerd://cd117befcc27e256954893238a6b2dc8787cce5200aa83e9f9d43f184de27ad0
//             exitCode: 0
//             finishedAt: "2022-06-09T17:20:57Z"
//             message: '[{"key":"commit","value":"b6f6ffabde2b1dca2d2ba498205793bc80619475","type":1},{"key":"url","value":"https://github.com/lilley2412/samples-go-api.git","type":1}]'
//             reason: Completed
//             startedAt: "2022-06-09T17:20:57Z"
//         taskResults:
//         - name: commit
//           value: b6f6ffabde2b1dca2d2ba498205793bc80619475
//         - name: url
//           value: https://github.com/lilley2412/samples-go-api.git
//         taskSpec:
//           description: |-
//             These Tasks are Git tasks to work with repositories used by other tasks in your Pipeline.
//             The git-clone Task will clone a repo from the provided url into the output Workspace. By default the repo will be cloned into the root of your Workspace. You can clone into a subdirectory by setting this Task's subdirectory param. This Task also supports sparse checkouts. To perform a sparse checkout, pass a list of comma separated directory patterns to this Task's sparseCheckoutDirectories param.
//           params:
//           - description: Repository URL to clone from.
//             name: url
//             type: string
//           - default: ""
//             description: Revision to checkout. (branch, tag, sha, ref, etc...)
//             name: revision
//             type: string
//           - default: ""
//             description: Refspec to fetch before checking out revision.
//             name: refspec
//             type: string
//           - default: "true"
//             description: Initialize and fetch git submodules.
//             name: submodules
//             type: string
//           - default: "1"
//             description: Perform a shallow clone, fetching only the most recent N
//               commits.
//             name: depth
//             type: string
//           - default: "true"
//             description: Set the ` + "`http.sslVerify`" + ` global git config. Setting this
//               to ` + "`false`" + ` is not advised unless you are sure that you trust your git
//               remote.
//             name: sslVerify
//             type: string
//           - default: ""
//             description: Subdirectory inside the ` + "`output`" + ` Workspace to clone the repo
//               into.
//             name: subdirectory
//             type: string
//           - default: ""
//             description: Define the directory patterns to match or exclude when performing
//               a sparse checkout.
//             name: sparseCheckoutDirectories
//             type: string
//           - default: "true"
//             description: Clean out the contents of the destination directory if it
//               already exists before cloning.
//             name: deleteExisting
//             type: string
//           - default: ""
//             description: HTTP proxy server for non-SSL requests.
//             name: httpProxy
//             type: string
//           - default: ""
//             description: HTTPS proxy server for SSL requests.
//             name: httpsProxy
//             type: string
//           - default: ""
//             description: Opt out of proxying HTTP/HTTPS requests.
//             name: noProxy
//             type: string
//           - default: "true"
//             description: Log the commands that are executed during ` + "`git-clone`" + `'s operation.
//             name: verbose
//             type: string
//           - default: gcr.io/tekton-releases/github.com/tektoncd/pipeline/cmd/git-init:v0.29.0
//             description: The image providing the git-init binary that this Task runs.
//             name: gitInitImage
//             type: string
//           - default: /tekton/home
//             description: |
//               Absolute path to the user's home directory. Set this explicitly if you are running the image as a non-root user or have overridden
//               the gitInitImage param with an image containing custom user configuration.
//             name: userHome
//             type: string
//           results:
//           - description: The precise commit SHA that was fetched by this Task.
//             name: commit
//           - description: The precise URL that was fetched by this Task.
//             name: url
//           steps:
//           - env:
//             - name: HOME
//               value: $(params.userHome)
//             - name: PARAM_URL
//               value: $(params.url)
//             - name: PARAM_REVISION
//               value: $(params.revision)
//             - name: PARAM_REFSPEC
//               value: $(params.refspec)
//             - name: PARAM_SUBMODULES
//               value: $(params.submodules)
//             - name: PARAM_DEPTH
//               value: $(params.depth)
//             - name: PARAM_SSL_VERIFY
//               value: $(params.sslVerify)
//             - name: PARAM_SUBDIRECTORY
//               value: $(params.subdirectory)
//             - name: PARAM_DELETE_EXISTING
//               value: $(params.deleteExisting)
//             - name: PARAM_HTTP_PROXY
//               value: $(params.httpProxy)
//             - name: PARAM_HTTPS_PROXY
//               value: $(params.httpsProxy)
//             - name: PARAM_NO_PROXY
//               value: $(params.noProxy)
//             - name: PARAM_VERBOSE
//               value: $(params.verbose)
//             - name: PARAM_SPARSE_CHECKOUT_DIRECTORIES
//               value: $(params.sparseCheckoutDirectories)
//             - name: PARAM_USER_HOME
//               value: $(params.userHome)
//             - name: WORKSPACE_OUTPUT_PATH
//               value: $(workspaces.output.path)
//             - name: WORKSPACE_SSH_DIRECTORY_BOUND
//               value: $(workspaces.ssh-directory.bound)
//             - name: WORKSPACE_SSH_DIRECTORY_PATH
//               value: $(workspaces.ssh-directory.path)
//             - name: WORKSPACE_BASIC_AUTH_DIRECTORY_BOUND
//               value: $(workspaces.basic-auth.bound)
//             - name: WORKSPACE_BASIC_AUTH_DIRECTORY_PATH
//               value: $(workspaces.basic-auth.path)
//             - name: WORKSPACE_SSL_CA_DIRECTORY_BOUND
//               value: $(workspaces.ssl-ca-directory.bound)
//             - name: WORKSPACE_SSL_CA_DIRECTORY_PATH
//               value: $(workspaces.ssl-ca-directory.path)
//             image: $(params.gitInitImage)
//             name: clone
//             resources: {}
//             script: |
//               #!/usr/bin/env sh
//               set -eu

//               if [ "${PARAM_VERBOSE}" = "true" ] ; then
//                 set -x
//               fi

//               if [ "${WORKSPACE_BASIC_AUTH_DIRECTORY_BOUND}" = "true" ] ; then
//                 cp "${WORKSPACE_BASIC_AUTH_DIRECTORY_PATH}/.git-credentials" "${PARAM_USER_HOME}/.git-credentials"
//                 cp "${WORKSPACE_BASIC_AUTH_DIRECTORY_PATH}/.gitconfig" "${PARAM_USER_HOME}/.gitconfig"
//                 chmod 400 "${PARAM_USER_HOME}/.git-credentials"
//                 chmod 400 "${PARAM_USER_HOME}/.gitconfig"
//               fi

//               if [ "${WORKSPACE_SSH_DIRECTORY_BOUND}" = "true" ] ; then
//                 cp -R "${WORKSPACE_SSH_DIRECTORY_PATH}" "${PARAM_USER_HOME}"/.ssh
//                 chmod 700 "${PARAM_USER_HOME}"/.ssh
//                 chmod -R 400 "${PARAM_USER_HOME}"/.ssh/*
//               fi

//               if [ "${WORKSPACE_SSL_CA_DIRECTORY_BOUND}" = "true" ] ; then
//                  export GIT_SSL_CAPATH="${WORKSPACE_SSL_CA_DIRECTORY_PATH}"
//               fi
//               CHECKOUT_DIR="${WORKSPACE_OUTPUT_PATH}/${PARAM_SUBDIRECTORY}"

//               cleandir() {
//                 # Delete any existing contents of the repo directory if it exists.
//                 #
//                 # We don't just "rm -rf ${CHECKOUT_DIR}" because ${CHECKOUT_DIR} might be "/"
//                 # or the root of a mounted volume.
//                 if [ -d "${CHECKOUT_DIR}" ] ; then
//                   # Delete non-hidden files and directories
//                   rm -rf "${CHECKOUT_DIR:?}"/*
//                   # Delete files and directories starting with . but excluding ..
//                   rm -rf "${CHECKOUT_DIR}"/.[!.]*
//                   # Delete files and directories starting with .. plus any other character
//                   rm -rf "${CHECKOUT_DIR}"/..?*
//                 fi
//               }

//               if [ "${PARAM_DELETE_EXISTING}" = "true" ] ; then
//                 cleandir
//               fi

//               test -z "${PARAM_HTTP_PROXY}" || export HTTP_PROXY="${PARAM_HTTP_PROXY}"
//               test -z "${PARAM_HTTPS_PROXY}" || export HTTPS_PROXY="${PARAM_HTTPS_PROXY}"
//               test -z "${PARAM_NO_PROXY}" || export NO_PROXY="${PARAM_NO_PROXY}"

//               /ko-app/git-init \
//                 -url="${PARAM_URL}" \
//                 -revision="${PARAM_REVISION}" \
//                 -refspec="${PARAM_REFSPEC}" \
//                 -path="${CHECKOUT_DIR}" \
//                 -sslVerify="${PARAM_SSL_VERIFY}" \
//                 -submodules="${PARAM_SUBMODULES}" \
//                 -depth="${PARAM_DEPTH}" \
//                 -sparseCheckoutDirectories="${PARAM_SPARSE_CHECKOUT_DIRECTORIES}"
//               cd "${CHECKOUT_DIR}"
//               RESULT_SHA="$(git rev-parse HEAD)"
//               EXIT_CODE="$?"
//               if [ "${EXIT_CODE}" != 0 ] ; then
//                 exit "${EXIT_CODE}"
//               fi
//               printf "%s" "${RESULT_SHA}" > "$(results.commit.path)"
//               printf "%s" "${PARAM_URL}" > "$(results.url.path)"
//           workspaces:
//           - description: The git repo will be cloned onto the volume backing this
//               Workspace.
//             name: output
//           - description: |
//               A .ssh directory with private key, known_hosts, config, etc. Copied to
//               the user's home before git commands are executed. Used to authenticate
//               with the git remote when performing the clone. Binding a Secret to this
//               Workspace is strongly recommended over other volume types.
//             name: ssh-directory
//             optional: true
//           - description: |
//               A Workspace containing a .gitconfig and .git-credentials file. These
//               will be copied to the user's home before any git commands are run. Any
//               other files in this Workspace are ignored. It is strongly recommended
//               to use ssh-directory over basic-auth whenever possible and to bind a
//               Secret to this Workspace over other volume types.
//             name: basic-auth
//             optional: true
//           - description: |
//               A workspace containing CA certificates, this will be used by Git to
//               verify the peer with when fetching or pushing over HTTPS.
//             name: ssl-ca-directory
//             optional: true
//     myapp-4zhw4-containerize:
//       pipelineTaskName: containerize
//       status:
//         completionTime: "2022-06-09T17:21:12Z"
//         conditions:
//         - lastTransitionTime: "2022-06-09T17:21:12Z"
//           message: All Steps have completed executing
//           reason: Succeeded
//           status: "True"
//           type: Succeeded
//         podName: myapp-4zhw4-containerize-pod
//         startTime: "2022-06-09T17:21:06Z"
//         steps:
//         - container: step-containerize
//           imageID: gcr.io/kaniko-project/executor@sha256:3bc3f3a05f803cac29164ce12617a7be64931748c944f6c419565f500b65e8db
//           name: containerize
//           terminated:
//             containerID: containerd://cc97f062214ead15d55ef07bee455703c2b28b8dac353b2b7d19a99c7b95857f
//             exitCode: 0
//             finishedAt: "2022-06-09T17:21:09Z"
//             reason: Completed
//             startedAt: "2022-06-09T17:21:09Z"
//         taskSpec:
//           steps:
//           - image: gcr.io/kaniko-project/executor:debug
//             name: containerize
//             resources: {}
//             script: |
//               #!/busybox/sh

//               /kaniko/executor --context=dir://$(workspaces.output.path) \
//                 --snapshotMode=redo \
//                 --insecure-registry 10.0.10.1:30500 \
//                 --dockerfile Dockerfile \
//                 --destination 10.0.10.1:30500/samples-go-api:latest
//           workspaces:
//           - name: output
// `
