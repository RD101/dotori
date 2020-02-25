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
$ yum -y install epel-release
$ rpm -Uvh http://li.nux.ro/download/nux/dextop/el7/x86_64/nux-dextop-release-0-5.el7.nux.noarch.rpm
$ yum install ffmpeg ffmpeg-devel -y
```

#### ocio.config 설치
원하는 컬러스페이스로 이미지의 컬러스페이스를 변환하기 위해서 OpenColorIO-Configs를 설치합니다.
```bash
$ yum install OpenColorIO
```

#### OpenImageIO
이미지를 컨버팅하기 위해서 OpenImageIO를 설치합니다.

```bash
$ yum install openimageio
```
