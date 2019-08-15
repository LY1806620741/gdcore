package gdcore
import (
    "errors"
    "reflect"
)
//class文件结构，首字母大小写来控制访问权
type ClassFile struct{
    Magic                []byte             //魔数
    Minor_version        int                //次要版本
    Major_version        int                //主要版本
    Constant_pool_count  int                //常量
    Constant_pool        []Cp_Info          //常量池，通过索引导向常量
    Access_flags         int                //访问权限
    This_class           int                //class文件自身索引
    Supe_class           int                //父类索引
    Interfaces_count     int                //实现的接口个数 ps:java一个类可以实现多个接口
    Interfaces           []int              //每个接口的索引
    Fields_count         int                //字段个数
    Fields               []Field_Info       //字段信息
    Methods_count        int                //方法数
    Methods              []Method_Info      //方法信息
    Attributes_count     int                //类属性数
    Attributes           []Attribute_Info   //类属性信息
}                        
//常量                   
type Cp_Info struct{     
    Tag                  int                 //常量标志(代表数据类型)
    Info                 []byte              //常量数据(各种类型的结构不一样，暂时不关心所有数据类型)
}                        
//字段                   
type Field_Info struct{  
    Access_flags         int                 //访问权限
    Name_index           int                 //名字索引，指向常量池UTF8类型
    Descriptor_index     int                 //描述符索引
    Attributes_count     int                 //字段属性的个数
    Attributes           []Attribute_Info    //字段属性详细
}                        
//方法                   
type Method_Info struct{ 
    Access_flags         int                 //权限
    Name_index           int                 //方法名常量池UTF8索引
    Descriptor_index     int                 //描述符索引
    Attributes_count     int                 //方法属性的个数
    Attributes           []Attribute_Info    //方法属性详细
}
//属性信息
type Attribute_Info struct{
    Attribute_name_index int                 //指向常量库utf8
    Attribute_length     int                 //长度
    Info                 []byte              //属性信息
}

//加载读取
func (c *ClassFile) Load(bs  []byte) (err error){//用非匿名返回值，defer可以拦截其返回
    //实现javap -v class类似功能，读取class文件到go数据结构中
    //检查魔数(class标志),0xCAFEBABE,3405691582,b'11001010111111101011101010111110',[]int{"202","254","186","190"}
    c.Magic=bs[:4]
    if c.Magic[0]==202&&c.Magic[1]==254&&c.Magic[2]==186&&c.Magic[3]==190{
        //常在岸边走，哪有不湿鞋，万一读取失败
        defer func() {
            if r := recover(); r != nil {
                err=errors.New("Load Fail, unknown class format -> "+reflect.ValueOf(r).String())
            }
        }()
        c.Minor_version=int(bs[4])*256+int(bs[5])                              //获取次要版本
        c.Major_version=int(bs[6])*256+int(bs[7])                              //获取主要版本
        c.Constant_pool_count=int(bs[8])*256+int(bs[9])                        //获取常量池大小
        //获取常量                                                             
        var readflag=10                                                        //读取位置标志位
        for i:=1;i<c.Constant_pool_count;i++{                                  //从一开始数
            cp := Cp_Info{}                                                    
            cp.Tag=int(bs[readflag])                                           
            switch cp.Tag{                                                     //对常量池进行分类，但不关心其余类型
                case 1:                                                        //CONSTANT_Utf8_info
                    var length=int(int(bs[readflag+1])*256+int(bs[readflag+2]))//长度信息读取后不需要
                    cp.Info=bs[readflag+3:readflag+3+length]
                    readflag+=3+length
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 3:                                                        //CONSTANT_Class_info format
                    cp.Info=bs[readflag+1:readflag+5]
                    readflag+=5
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 4:                                                        //CONSTANT_Float_info
                    cp.Info=bs[readflag+1:readflag+5]
                    readflag+=5
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 5:                                                        //CONSTANT_Long_info  All 8-byte constants take up two entries in the constant_pool table of the class file. If a CONSTANT_Long_info or CONSTANT_Double_info structure is the item in the constant_pool table at index n, then the next usable item in the pool is located at index n+2. The constant_pool index n+1 must be valid but is considered unusable.
                    cp.Info=bs[readflag+1:readflag+9]
                    readflag+=9
                    c.Constant_pool=append(c.Constant_pool,cp)
                    i+=1
                case 6:                                                        //CONSTANT_Double_info
                    cp.Info=bs[readflag+1:readflag+9]
                    readflag+=9
                    c.Constant_pool=append(c.Constant_pool,cp)
                    i+=1
                case 7:                                                        //CONSTANT_Class_info format
                    cp.Info=bs[readflag+1:readflag+3]
                    readflag+=3
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 8:                                                        //CONSTANT_String_info
                    cp.Info=bs[readflag+1:readflag+3]
                    readflag+=3
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 9:                                                        //CONSTANT_Fieldref_info
                    cp.Info=bs[readflag+1:readflag+5]
                    readflag+=5
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 10:                                                        //CONSTANT_Methodref_info
                    cp.Info=bs[readflag+1:readflag+5]
                    readflag+=5
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 11:                                                        //CONSTANT_InterfaceMethodref_info
                    cp.Info=bs[readflag+1:readflag+5]
                    readflag+=5
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 12:                                                        //CONSTANT_NameAndType_info
                    cp.Info=bs[readflag+1:readflag+5]
                    readflag+=5
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 15:                                                        //CONSTANT_MethodHandle_info
                    cp.Info=bs[readflag+1:readflag+5]
                    readflag+=5
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 16:                                                        //CONSTANT_MethodType_info
                    cp.Info=bs[readflag+1:readflag+3]
                    readflag+=3
                    c.Constant_pool=append(c.Constant_pool,cp)
                case 18:                                                        //CONSTANT_InvokeDynamic_info
                    cp.Info=bs[readflag+1:readflag+5]
                    readflag+=5
                    c.Constant_pool=append(c.Constant_pool,cp)
                default:                                                        //未知的跳过
                    i-=1
                    readflag+=1
            }
        }
        //读取指定类或接口的访问权限
        c.Access_flags = int(bs[readflag])*256+int(bs[readflag+1])
        readflag+=2
        //读取class文件的索引
        c.This_class=int(bs[readflag])*256+int(bs[readflag+1])
        readflag+=2
        //读取父类的索引
        c.Supe_class=int(bs[readflag])*256+int(bs[readflag+1])
        readflag+=2
        //读取当前class实现的接口
        c.Interfaces_count=int(bs[readflag])*256+int(bs[readflag+1])
        readflag+=2
        for i:=0;i<c.Interfaces_count;i++{
            c.Interfaces=append(c.Interfaces,int(bs[readflag])*256+int(bs[readflag+1]))
            readflag+=2
        }
        //读取class中定义的字段
        c.Fields_count=int(bs[readflag])*256+int(bs[readflag+1])
        readflag+=2
        for i:=0;i<c.Fields_count;i++{
            f := Field_Info{}
            f.Access_flags=int(bs[readflag])*256+int(bs[readflag+1])
            f.Name_index=int(bs[readflag+2])*256+int(bs[readflag+3])
            f.Descriptor_index=int(bs[readflag+4])*256+int(bs[readflag+5])
            f.Attributes_count=int(bs[readflag+6])*256+int(bs[readflag+7])
            readflag+=8
            //读取字段中定义的属性
            for j:=0;j<f.Attributes_count;j++{
                a := Attribute_Info{}
                a.Attribute_name_index = int(bs[readflag])*256+int(bs[readflag+1])
                readflag+=2
                var tmp=int(bs[readflag])*256+int(bs[readflag+1])
                readflag+=2
                a.Attribute_length = tmp*256*256+int(bs[readflag])*256+int(bs[readflag+1])
                readflag+=2
                a.Info=bs[readflag:readflag+a.Attribute_length]
                readflag+=a.Attribute_length
                f.Attributes=append(f.Attributes,a)
            }
            c.Fields=append(c.Fields,f)
        }
        //读取class中定义的方法
        c.Methods_count=int(bs[readflag])*256+int(bs[readflag+1])
        readflag+=2
        for i:=0;i<c.Methods_count;i++{
            m := Method_Info{}
            m.Access_flags=int(bs[readflag])*256+int(bs[readflag+1])
            m.Name_index=int(bs[readflag+2])*256+int(bs[readflag+3])
            m.Descriptor_index=int(bs[readflag+4])*256+int(bs[readflag+5])
            m.Attributes_count=int(bs[readflag+6])*256+int(bs[readflag+7])
            readflag+=8
            //读取方法中定义的属性
            for j:=0;j<m.Attributes_count;j++{
                a := Attribute_Info{}
                a.Attribute_name_index = int(bs[readflag])*256+int(bs[readflag+1])
                readflag+=2
                var tmp=int(bs[readflag])*256+int(bs[readflag+1])
                readflag+=2
                a.Attribute_length = tmp*256*256+int(bs[readflag])*256+int(bs[readflag+1])
                readflag+=2
                a.Info=bs[readflag:readflag+a.Attribute_length]
                readflag+=a.Attribute_length
                m.Attributes=append(m.Attributes,a)
            }
            c.Methods=append(c.Methods,m)
        }
        //读取class中定义的属性
        c.Attributes_count=int(bs[readflag])*256+int(bs[readflag+1])
        readflag+=2
        for j:=0;j<c.Attributes_count;j++{
                a := Attribute_Info{}
                a.Attribute_name_index = int(bs[readflag])*256+int(bs[readflag+1])
                //fmt.Println(a.Attribute_name_index)
                readflag+=2
                var tmp=int(bs[readflag])*256+int(bs[readflag+1])
                readflag+=2
                a.Attribute_length = tmp*256*256+int(bs[readflag])*256+int(bs[readflag+1])
                //fmt.Println(a.Attribute_length)
                readflag+=2
                a.Info=bs[readflag:readflag+a.Attribute_length]
                readflag+=a.Attribute_length
                c.Attributes=append(c.Attributes,a)
            }
        //检查是否剩余byte
        if (readflag!= len(bs)){
            return errors.New("Error. Unknown java class")
        }
        //输出剩余byte(调试使用)
        //fmt.Printf("%s\n%d\n%X\n",bs[readflag:],bs[readflag:],bs[readflag:])
        //fmt.Println(c)
        return nil
    }else{
        return errors.New("Error. This isn't java class")
    }
}