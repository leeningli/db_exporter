需要下载的：<br />
"github.com/wendal/go-oci8" <br />
oracle-instantclient11.2-basic-11.2.0.1.0-1.x86_64.zip<br />
oracle-instantclient11.2-sdk-11.2.0.1.0-1.x86_64.zip<br />
pkgconfig-0.17.2.tar.bz2<br />

unzip oracle-instantclient11.2-basic-11.2.0.1.0-1.x86_64.zip<br />
unzip oracle-instantclient11.2-sdk-11.2.0.1.0-1.x86_64.zip
tar -zcvf instantclient_11_2.tgz instantclient_11_2
mv instantclient_11_2.tgz /usr/lib
cd /usr/lib
tar -zxvf instantclient_11_2.tgz
ln /usr/lib/instantclient_11_2/libclntsh.so.11.1 /usr/lib/libclntsh.so
ln /usr/lib/instantclient_11_2/libocci.so.11.1 /usr/lib/libocci.so
ln /usr/lib/instantclient_11_2/libociei.so /usr/lib/libociei.so
ln /usr/lib/instantclient_11_2/libnnz11.so /usr/lib/libnnz11.so

tar -xvf pkgconfig-0.17.2.tar.bz2
cd ;./configure make make install
/usr/lib64/pkgconfig
vim oci8.pc
	prefix=/usr/lib/instantclient_11_2 
	libdir=${prefix}
	includedir=${prefix}/sdk/include
	Name: OCI
	Description: Oracle database engine
	Version: 11.2                                            
	Libs: -L${libdir} -lclntsh
	Libs.private: 
	Cflags: -I${includedir}

7、.bashrc 文件中添加系统变量
export GOHOME=/docker/home/docker/go
export GOROOT=$GOHOME/go
export GOPATH=$GOHOME/myproject
export ORACLE_HOME=$GOHOME/instantclient_11_2
export TNS_ADMIN=$ORACLE_HOME/network/admin
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$ORACLE_HOME
export PKG_CONFIG_PATH=/usr/lib64/pkgconfig
PATH=$PATH:$HOME/bin:$GOROOT/bin

$GOHOME/instantclient_11_2下创建network/admin，将tnsnames.ora文件放置此处

119行：(**C.OCIServer)(unsafe.Pointer(&conn.svc)),---》(**C.OCISvcCtx)(unsafe.Pointer(&conn.svc)),
136行：(*C.OCIServer)(c.svc),---》(*C.OCISvcCtx)(c.svc),
263行：(*C.OCIServer)(c.svc), (*C.OCISvcCtx)(s.c.svc),
383行：(*C.OCIServer)(c.svc), (*C.OCISvcCtx)(s.c.svc),
