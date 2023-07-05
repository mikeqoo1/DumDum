NICI=NiciServer.out

ifeq ($(google),6016)
IsGoogle=NO
endif

ifeq ($(pvc),on)
IsPvc=YES
endif

.PHONY: build clean install help

build:
	go build -ldflags="-X main.IsGoogle=${IsGoogle} -X main.IsPvc=${IsPvc}" -o bin/${NICI} main.go

install:
	go install

clean:
	if [ -f bin/${NICI} ] ; then rm bin/${NICI} ; fi

help:
	@echo "make 格式化"
	@echo "make build 編譯程式碼產生執行檔"
	@echo "make clean 清除執行檔"
	@echo "make test 執行單元測試"
	@echo "make check 格式化go程式碼"
	@echo "make cover 檢查測試程式碼的覆蓋率"
	@echo "make run 直接跑程式"
	@echo "make lint 程式碼檢查"
	@echo "make docker 建構docker image"