package gform

import (
    "syscall"
    "unsafe"
    "w32"
    "w32/user32"
    "w32/comdlg32"
    "w32/ole32"
    "w32/shell32"
)

func genOFN(parent Controller, title, filter string, filterIndex uint, initialDir string, buf []uint16) *w32.OPENFILENAME {
    var ofn w32.OPENFILENAME
    ofn.StructSize = uint(unsafe.Sizeof(ofn))
    ofn.Owner = parent.Handle()

    if filter != "" {
        filterBuf := make([]uint16, len(filter)+1)
        copy(filterBuf, syscall.StringToUTF16(filter))
        // Replace '|' with the expcted '\0'
        for i, c := range filterBuf {
            if byte(c) == '|' {
                filterBuf[i] = uint16(0)
            }
        }
        ofn.Filter = &filterBuf[0]
        ofn.FilterIndex = uint(filterIndex)
    }

    ofn.File = &buf[0]
    ofn.MaxFile = uint(len(buf))

    if initialDir != "" {
        ofn.InitialDir = syscall.StringToUTF16Ptr(initialDir)
    }
    if title != "" {
        ofn.Title = syscall.StringToUTF16Ptr(title)
    }

    ofn.Flags = w32.OFN_FILEMUSTEXIST

    return &ofn
}

func ShowOpenFileDlg(parent Controller, title, filter string, filterIndex uint, initialDir string) (filePath string, accepted bool) {
    buf := make([]uint16, 1024)
    ofn := genOFN(parent, title, filter, filterIndex, initialDir, buf)

    if accepted = comdlg32.GetOpenFileName(ofn); accepted {
        filePath = syscall.UTF16ToString(buf)
    }

    return
}

func ShowSaveFileDlg(parent Controller, title, filter string, filterIndex uint, initialDir string) (filePath string, accepted bool) {
    buf := make([]uint16, 1024)
    ofn := genOFN(parent, title, filter, filterIndex, initialDir, buf)

    if accepted = comdlg32.GetSaveFileName(ofn); accepted {
        filePath = syscall.UTF16ToString(buf)
    }

    return
}

func ShowBrowseFolderDlg(parent Controller, title string) (folder string, accepted bool) {
    var bi w32.BROWSEINFO
    bi.Owner = parent.Handle()
    bi.Title = syscall.StringToUTF16Ptr(title)
    bi.Flags = w32.BIF_RETURNONLYFSDIRS | w32.BIF_NEWDIALOGSTYLE

    ole32.CoInitialize()
    ret := shell32.SHBrowseForFolder(&bi)
    ole32.CoUninitialize()

    folder = shell32.SHGetPathFromIDList(ret)
    accepted = folder != ""
    return
}

func MsgBox(parent Controller, title, caption string, flags uint) int {
    var result int
    if parent != nil {
        result = user32.MessageBox(parent.Handle(), caption, title, flags)
    } else {
        result = user32.MessageBox(0, caption, title, flags)
    }

    return result
}