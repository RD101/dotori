# Library Setting

## macOS
brew가 필요합니다. brew를 먼저 설치해주세요.

#### FFmpeg
동영상을 변환하기 위해서 FFmpeg를 설치합니다.

```bash
$ brew install ffmpeg
```

#### ocio.config 설치
원하는 컬러스페이스로 이미지의 컬러스페이스를 변환하기 위해서 OpenColorIO-Configs를 설치합니다.

```bash
$ cd ~
$ git clone https://github.com/imageworks/OpenColorIO-Configs
```

#### OpenImageIO
이미지를 컨버팅하기 위해서 OpenImageIO를 설치합니다.

```bash
$ brew install openimageio
```

## CentOS
yum 으로 설치를 합니다.

#### FFmpeg
동영상을 변환하기 위해서 FFmpeg를 설치합니다.
```bash
$ cd ~
$ mkdir -p app/ffmpeg
$ cd app/ffmpeg/
$ wget http://johnvansickle.com/ffmpeg/builds/ffmpeg-git-amd64-static.tar.xz
$ tar xpvf ffmpeg-git-amd64-static.tar.xz --strip 1
```

#### ocio.config 설치
원하는 컬러스페이스로 이미지의 컬러스페이스를 변환하기 위해서 OpenColorIO-Configs를 설치합니다.
```bash
$ cd ~
$ sudo yum install git // CentOS를 최초 설치하면 Git이 설치되어있지 않다. 설치한다.
$ git clone https://github.com/imageworks/OpenColorIO-Configs
```

#### OpenImageIO
이미지를 컨버팅하기 위해서 OpenImageIO를 설치합니다.

```bash
$ yum install OpenImageIO
$ yum install OpenImageIO-iv
$ yum install OpenImageIO-devel
$ yum install OpenImageIO-utils
$ yum install python-OpenImageIO
```
