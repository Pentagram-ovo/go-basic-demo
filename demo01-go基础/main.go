package main

import (
	"fmt"
	"os"
)

type Student struct {
	Name   string
	Age    int
	Id     string
	Gender string
}

var studentMap = make(map[string]Student)

func AddStudent(s Student) error {
	if _, ok := studentMap[s.Id]; ok {
		return fmt.Errorf("此学号已存在")
	}
	studentMap[s.Id] = s
	return nil
}

func DeleteStudent(id string) error {
	if _, ok := studentMap[id]; !ok {
		return fmt.Errorf("该学号不存在")
	}
	delete(studentMap, id)
	return nil
}

func UpdateStudent(s Student, id string) error {
	if _, ok := studentMap[id]; !ok {
		return fmt.Errorf("该学号不存在")
	}
	studentMap[id] = s
	return nil
}

func ShowStudent() {
	if len(studentMap) == 0 {
		fmt.Println("暂无学生！请添加~")
		return
	}
	fmt.Println("全部学生信息如下:")
	num := 0
	for _, v := range studentMap {
		num++
		fmt.Printf("第%d位学生的学号为：%s 姓名为：%s 性别为：%s 年龄为：%d \n", num, v.Id, v.Name, v.Gender, v.Age)
	}
}

func main() {
	for {
	A:
		fmt.Println("欢迎来到五角星的学生管理系统~~~")
		fmt.Println("1.查询学生(按学号)")
		fmt.Println("2.新增学生")
		fmt.Println("3.删除学生")
		fmt.Println("4.修改学生信息")
		fmt.Println("5.显示全部学生")
		fmt.Println("6.退出系统")
		var num int
		fmt.Println("请输入数字以实现对于的功能：")
		fmt.Scanf("%d", &num)
		switch num {
		case 1:
			fmt.Println("请输入要查询学生的学号：")
			var id string
			fmt.Scanf("%s", &id)
			if _, ok := studentMap[id]; !ok {
				fmt.Println("学生不存在！")
			} else {
				fmt.Println("查询成功！")
				fmt.Printf("学生的学号为：%s 姓名为：%s 性别为：%s 年龄为：%d \n", studentMap[id].Id, studentMap[id].Name, studentMap[id].Gender, studentMap[id].Age)
			}
			fmt.Println("按回车继续...")
			fmt.Scanln()
			goto A
		case 2:
			var s Student
			fmt.Println("请输入要添加的学生的信息：")
			var id string
			fmt.Println("请输入学生学号：")
			fmt.Scanf("%s", &id)
			s.Id = id
			var gender string
			fmt.Println("请输入学生性别：")
			fmt.Scanf("%s", &gender)
			s.Gender = gender
			var age int
			fmt.Println("请输入学生年龄：")
			fmt.Scanf("%d", &age)
			s.Age = age
			var name string
			fmt.Println("请输入学生姓名：")
			fmt.Scanf("%s", &name)
			s.Name = name
			error := AddStudent(s)
			if error != nil {
				fmt.Println("操作失败！")
			} else {
				fmt.Println("操作成功！")
			}
			fmt.Println("按回车继续...")
			fmt.Scanln()
			goto A
		case 3:
			fmt.Println("请输入要删除的学生的学号：")
			var id string
			fmt.Scanf("%s", &id)
			error := DeleteStudent(id)
			if error != nil {
				fmt.Println("操作失败！")
			} else {
				fmt.Println("操作成功！")
			}
			fmt.Println("按回车继续...")
			fmt.Scanln()
			goto A
		case 4:
			var s Student
			fmt.Println("请输入要修改的学生的信息：")
			var id_old string
			fmt.Println("请输入修改前学生学号：")
			fmt.Scanf("%s", &id_old)
			var id_new string
			fmt.Println("请输入修改后学生学号：")
			fmt.Scanf("%s", &id_new)
			s.Id = id_new
			var gender string
			fmt.Println("请输入学生性别：")
			fmt.Scanf("%s", &gender)
			s.Gender = gender
			var age int
			fmt.Println("请输入学生年龄：")
			fmt.Scanf("%d", &age)
			s.Age = age
			var name string
			fmt.Println("请输入学生姓名：")
			fmt.Scanf("%s", &name)
			s.Name = name
			error := UpdateStudent(s, id_old)
			if error != nil {
				fmt.Println("操作失败！")
			} else {
				fmt.Println("操作成功！")
			}
			fmt.Println("按回车继续...")
			fmt.Scanln()
			goto A
		case 5:
			ShowStudent()
			fmt.Println("按回车继续...")
			fmt.Scanln()
			goto A
		case 6:
			fmt.Println("谢谢使用，欢迎下次操作 ovo~")
			os.Exit(0)
		default:
			fmt.Println("无效的操作，请重新输入")
			fmt.Println("按回车继续...")
			fmt.Scanln()
			goto A
		}
	}
}
