#db_exporter

install oracle-driver on centos/rhel
1.first:you hava get these file file as follow:
"github.com/wendal/go-oci8"
oracle-instantclient11.2-basic-11.2.0.1.0-1.x86_64.zip
oracle-instantclient11.2-sdk-11.2.0.1.0-1.x86_64.zip
pkgconfig-0.17.2.tar.bz2

2.second,install these
unzip oracle-instantclient11.2-basic-11.2.0.1.0-1.x86_64.zip
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

3:set your env
export GOHOME=/docker/home/docker/go
export GOROOT=$GOHOME/go
export GOPATH=$GOHOME/myproject
export ORACLE_HOME=$GOHOME/instantclient_11_2
export TNS_ADMIN=$ORACLE_HOME/network/admin
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$ORACLE_HOME
export PKG_CONFIG_PATH=/usr/lib64/pkgconfig
PATH=$PATH:$HOME/bin:$GOROOT/bin

4.create tnsnames.ora
create dir:network/admin on $GOHOME/instantclient_11_2
create file:tnsnames.ora on $GOHOME/instantclient_11_2/network/admin

5.change the source code:
vim /src/github.com/wendal/go-oci8/oci8.go
119:(**C.OCIServer)(unsafe.Pointer(&conn.svc)),----->(**C.OCISvcCtx)(unsafe.Pointer(&conn.svc)),
136:(*C.OCIServer)(c.svc),----->(*C.OCISvcCtx)(c.svc),
263:(*C.OCIServer)(c.svc),----->(*C.OCISvcCtx)(s.c.svc),
383:(*C.OCIServer)(c.svc),----->(*C.OCISvcCtx)(s.c.svc),