echo 'export M_VERSION="dev"' >> ~/.bashrc
echo 'export M_WORKDIR="$(pwd)/workdir"' >> ~/.bashrc
echo 'export M_RESOURCES="$(pwd)/resources"' >> ~/.bashrc
echo 'export M_SHARED="$(pwd)/shared"' >> ~/.bashrc
cd $(pwd)/resources/terraform && terraform init
