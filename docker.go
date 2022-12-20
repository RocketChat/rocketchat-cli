package main

//func CreateDockerErrorUIPage() {
//	layout := tview.NewFlex().SetDirection(tview.FlexRow)
//
//	errorMessage := `
//This setup demands a local docker installation. Please install docker and try again:
//
//https://docs.docker.com/engine/install/
//`
//
//	// Error
//	uiError := tview.NewTextView().
//		SetDynamicColors(true).
//		SetRegions(true).
//		SetChangedFunc(func() {
//			app.Draw()
//		})
//	uiError.SetText(errorMessage)
//
//	layout.AddItem(tview.NewFlex().AddItem(uiError, 0, 1, false), 5, 1, false)
//
//	// Menu Items
//	uiMenuList := tview.NewList().
//		AddItem("Quit", "", 'q', func() {
//			app.Stop()
//		})
//
//	layout.AddItem(tview.NewFlex().
//		AddItem(uiMenuList, 0, 1, true), 0, 1, true)
//
//	// Add page
//	cliPages.AddPage("Docker Error", layout, true, false)
//}
//
//func DockerInstalled() bool {
//	// Check if docker is installed
//	_, err := exec.Command("docker", "info").Output()
//
//	return err == nil
//}
