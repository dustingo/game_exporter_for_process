#### game_exporter

- version:1.0

- 执行:

  ```shell
  cd game_exporter
  ./game_exporter
  #帮助
  ./game_exporter -h
  
  ```

- 注意：

  ```text
  配置文件固定，在当前执行目录下，名字为：gameprocess.yml
  ```

- 配置文件说明

  ```text
  1.name: 进程名字，一般为缩写【唯一】
  2.cmdline: 根据cmdline去查找进程，为了更加准确或者区分特殊情况，最多可以有两个元素
  ```

  