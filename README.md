需要下载的：<br />
"github.com/wendal/go-oci8" <br />
oracle-instantclient11.2-basic-11.2.0.1.0-1.x86_64.zip<br />
oracle-instantclient11.2-sdk-11.2.0.1.0-1.x86_64.zip<br />
pkgconfig-0.17.2.tar.bz2<br />

unzip oracle-instantclient11.2-basic-11.2.0.1.0-1.x86_64.zip<br />
unzip oracle-instantclient11.2-sdk-11.2.0.1.0-1.x86_64.zip<br />
tar -zcvf instantclient_11_2.tgz instantclient_11_2<br />
mv instantclient_11_2.tgz /usr/lib<br />
cd /usr/lib<br />
tar -zxvf instantclient_11_2.tgz<br />
ln /usr/lib/instantclient_11_2/libclntsh.so.11.1 /usr/lib/libclntsh.so<br />
ln /usr/lib/instantclient_11_2/libocci.so.11.1 /usr/lib/libocci.so<br />
ln /usr/lib/instantclient_11_2/libociei.so /usr/lib/libociei.so<br />
ln /usr/lib/instantclient_11_2/libnnz11.so /usr/lib/libnnz11.so<br />

tar -xvf pkgconfig-0.17.2.tar.bz2<br />
cd ;./configure make make install<br />
/usr/lib64/pkgconfig<br />
vim oci8.pc<br />
	prefix=/usr/lib/instantclient_11_2 <br />
	libdir=${prefix}<br />
	includedir=${prefix}/sdk/include<br />
	Name: OCI<br />
	Description: Oracle database engine<br />
	Version: 11.2                                            <br />
	Libs: -L${libdir} -lclntsh<br />
	Libs.private: <br />
	Cflags: -I${includedir}<br />

7、.bashrc 文件中添加系统变量<br />
export GOHOME=/docker/home/docker/go<br />
export GOROOT=$GOHOME/go<br />
export GOPATH=$GOHOME/myproject<br />
export ORACLE_HOME=$GOHOME/instantclient_11_2<br />
export TNS_ADMIN=$ORACLE_HOME/network/admin<br />
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$ORACLE_HOME<br />
export PKG_CONFIG_PATH=/usr/lib64/pkgconfig<br />
PATH=$PATH:$HOME/bin:$GOROOT/bin<br />

$GOHOME/instantclient_11_2下创建network/admin，将tnsnames.ora文件放置此处<br />

119行：(**C.OCIServer)(unsafe.Pointer(&conn.svc)),---》(**C.OCISvcCtx)(unsafe.Pointer(&conn.svc)),<br />
136行：(*C.OCIServer)(c.svc),---》(*C.OCISvcCtx)(c.svc),<br />
263行：(*C.OCIServer)(c.svc), (*C.OCISvcCtx)(s.c.svc),<br />
383行：(*C.OCIServer)(c.svc), (*C.OCISvcCtx)(s.c.svc),<br />
