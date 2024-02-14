# Python Worker container build

#### Running a Dataplane Container

1.  **Run Dataplane Container:** To start a Dataplane container on your local machine, use the following command:

    ```sh
    docker run -p 9001:9000 --rm dataplane/dataplane:0.0.2
    ```

    This command does the following:

    * `-p 9001:9000` maps port 9001 on the host to port 9000 on the container, which is the default port for Dataplane.
    * `--rm` automatically removes the container when it exits.
    * `dataplane/dataplane:0.0.2` specifies the image to use, which is version 0.0.2 of Dataplane.

#### Building and Pushing Images to Docker Hub

2.  **Login to Docker Hub:** Before pushing images, you need to log in to Docker Hub:

    ```sh
    docker login
    ```

    Enter your Docker Hub username and password when prompted.
3.  **Set Dataplane Version Environment Variable:** Set an environment variable for the Dataplane version you're building:

    ```sh
    export dpversion=0.0.x
    ```

    Replace `0.0.x` with the actual version number you're working with.
4.  **Build Dataplane Main Image:** Navigate to the root directory of the Dataplane repository and build the main Dataplane image:

    ```sh
    docker build -t dataplane/dataplane:$dpversion -f docker-build/Dockerfile.main.alpine .
    ```
5.  **Tag and Push Main Image:** Tag the built image with the version number and push it to Docker Hub:

    ```sh
    docker push dataplane/dataplane:$dpversion
    ```

    Tag the image as `latest` and push it:

    ```sh
    docker tag dataplane/dataplane:$dpversion dataplane/dataplane:latest
    docker push dataplane/dataplane:latest
    ```
6.  **Build Dataplane Python Worker Image (Debian):** Build the Python worker image based on Debian:

    ```sh
    docker build -t dataplane/dataplane-worker-python:$dpversion -f docker-build/Dockerfile.workerpython.debian .
    ```
7.  **Tag and Push Python Worker Image (Debian):** Tag and push the Debian-based Python worker image:

    ```sh
    docker push dataplane/dataplane-worker-python:$dpversion
    ```

    Tag the image as `latest` and push it:

    ```sh
    docker tag dataplane/dataplane-worker-python:$dpversion dataplane/dataplane-worker-python:latest
    docker push dataplane/dataplane-worker-python:latest
    ```
8.  **Build Dataplane Python Worker Image (Ubuntu):** Build the Python worker image based on Ubuntu:

    ```sh
    docker build -t dataplane/dataplane-worker-python-ubuntu:latest -f docker-build/Dockerfile.workerpython.ubuntu .
    ```
9.  **Tag and Push Python Worker Image (Ubuntu):** Since we're using `latest` in the build step, we only need to push the Ubuntu-based Python worker image:

    ```sh
    docker push dataplane/dataplane-worker-python-ubuntu:latest
    ```

#### Notes:

* Ensure you have the necessary permissions to push images to the `dataplane` Docker Hub repository.
* Always replace the `$dpversion` placeholder with the actual version number you intend to use.
* The `latest` tag is typically used for the most recent stable version of an image.
* The `.` at the end of the `docker build` command indicates that Docker should use the current directory as the build context.
* Make sure you are in the root directory of the Dataplane repository when executing the build commands, as indicated by the comment `# NB: From root directory of this repo (not this directory):`.

By following these steps, you should be able to build and push Dataplane and its worker images to Docker Hub successfully.
