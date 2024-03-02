# [Prerequisite]
# - You must be using WSL
# - NSIS must be installed.
# - Go must be installed

ARCH = x86_64

SED = sed
GO = go
EXE = RemoteCmdExec.exe
SRC = RemoteCmdExec.go
BUILD_OPTIONS = GOOS=windows GOARCH=amd64

.PHONE: all
all: $(EXE)

$(EXE): $(SRC)
	$(BUILD_OPTIONS) $(GO) build -o $@ $^

.PHONE: clean
clean:
	-rm $(EXE) $(NSI_DST) $(INSTALLER)
