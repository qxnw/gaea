#########################################################################
# File Name: build_project.sh
# Author: ChamPly
# mail: champly@outlook.com
# Created Time: Tue Jul 11 16:46:12 2017
#########################################################################
#!/bin/bash

#项目路径
hydra_path=github.com/qxnw/hydra
project_path=github.com/qxnw/hydra2_core
# #项目名称

cd $GOPATH/src/${hydra_path}
git pull origin master

cd $GOPATH/src/${project_path}
git pull origin master

cd $GOPATH/bin
go install ${hydra_path}

go build -buildmode=plugin ${project_path}

tar zcvf hydra hydra2_core.tar.gz 
tar zcvf hydra2_core.so hydra2_core.tar.gz 

mv hydra2_core.tar.gz  ../3.0.2.1/hydra2_core.tar.gz


# echo "编译文件……"
# go build ${project_folder1}
# go build -buildmode=plugin  ${project_folder2}



# echo "编译完成，开始打包上传文件"
# tar -cf - ${project_name1} | ssh root@${server_ip} "cd ${server_folder_name}; rm -rf ${server_folder_name}/${project_name1}.bak;mv ${server_folder_name}/${project_name1} ${server_folder_name}/${project_name1}.bak; tar -xf -"
# tar -cf - ${project_name2} | ssh root@${server_ip} "cd ${server_folder_name}; rm -rf ${server_folder_name}/${project_name2}.bak;mv ${server_folder_name}/${project_name2} ${server_folder_name}/${project_name2}.bak; tar -xf -"
