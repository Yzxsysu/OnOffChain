echo "Shell 传递参数实例！";
echo "执行的文件名：$0";
echo "第一个参数为：$1";


# shellcheck disable=SC2077
if [ $1='f' ]
  then
  echo "第一个参数为：$1";
fi