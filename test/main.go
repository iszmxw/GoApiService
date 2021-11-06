package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"time"
)

func main() {
	var trip Trip
	trip = Trip{Icon: make([]int,4)}
	for{
		var a int
		fmt.Println("请选择操作 1 获取图标截图，2 点击图标 3 锁定聊天窗口并发送内容")
		fmt.Scanln(&a)
		switch a {
		case 1:
			trip.getScreen()
			trip.getShot()
		case 2:
			trip.clickIcon()
		case 3:
			trip.clickMessage()
		}
	}
}

type Trip struct {
	Icon []int
}
func (t *Trip)getScreen()  {
	fmt.Println("请在要打开的程序图标左上角点击左键")
	ok := robotgo.AddEvent("mleft")
	if ok{
		t.Icon[0],t.Icon[1] = robotgo.GetMousePos()
		fmt.Println("---请在要截图右下角下左键---",t.Icon)
	}
	ok = robotgo.AddEvent("mleft")
	if ok{
		x,y := robotgo.GetMousePos()
		t.Icon[2] = x - t.Icon[0]
		t.Icon[3] = y - t.Icon[1]
	}

}
func (t *Trip)getShot(){  // 获取程序图标截图
	fmt.Println("正在获取截图....")//392 678
	bit_map := robotgo.CaptureScreen(t.Icon...)
	robotgo.SaveBitmap(bit_map, "./appicon.png")//保存位图为图片
	fmt.Println("已保存截图")
}
func (t *Trip)clickIcon(){
	x,y := t.getImgXY("./appicon.png")
	if x < 0 || y < 0{
		fmt.Println("获取程序坐标失败")
		return
	}
	t.clickShot(x,y)
	t.clickMessage()
}
func (t *Trip)clickShot(x,y int)  {
	robotgo.MouseToggle("up")
	//robotgo.MoveMouseSmooth(x,y)
	//robotgo.MouseClick()
	robotgo.MoveClick(x,y,"left",true)
}
func (t *Trip)getImgXY(img string)(x,y int)  {
	bit_map := robotgo.OpenBitmap(img)
	fx, fy := robotgo.FindBitmap(bit_map)
	fmt.Println(fx,"--",fy)
	if fx < 0 && fy < 0{
		fmt.Println("获取坐标失败，重试获取")
		time.Sleep(time.Second*1)
		t.getImgXY(img)
	}
	return fx,fy
}
func (t *Trip)clickMessage()  {  // 锁定聊天窗口输入框
	x,y := t.getImgXY("./img/sendwindow.png")
	if x < 0 || y < 0{
		fmt.Println("获取程序坐标失败")
		return
	}
	t.clickShot(x,y)
	t.sendMessage("Hello Word!")
}
func (t *Trip)sendMessage(msg string)  {
	time.Sleep(time.Second)
	robotgo.TypeString(msg)
	time.Sleep(time.Second)
	robotgo.KeyTap("enter")
	time.Sleep(time.Second)
	robotgo.KeyTap("lctrl","enter")
	t.sendOut()
}
func(t *Trip)sendOut(){
	x,y := t.getImgXY("./img/sendbuttun.png")
	if x < 0 || y < 0{
		fmt.Println("获取程序坐标失败")
		return
	}
	t.clickShot(x,y)
}
