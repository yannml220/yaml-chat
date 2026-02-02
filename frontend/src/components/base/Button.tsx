import styled from '@emotion/styled'
import {forwardRef} from 'react';


interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement>, React.HTMLAttributes<HTMLButtonElement>{
    role? : 'primary' | 'primaryLight' | 'danger' | 'success' |  'warning' | 'secondary' | 'tertiary' ;
    variant? : 'outlined' | 'filled' | 'text' | 'gradient' ;
    colorVariant? : 'dark' |  'light' | 'medium' ;
    onClickHandler? : (e:React.MouseEvent<HTMLButtonElement>)=>void ; 
    color? : string ;
    bgColor? : string ;
	isActive? : boolean ;
	disabled? : boolean ;
    activeBgColor? : string ;
    activeColor? : string ;
    children? : React.ReactNode ;   
    text? : string ;
    padding? : string  ;
    borderRadius? : string ; 
    hoverBgColor? : string ;
    hoverColor? : string ;
    hoverBorder?:string ;
	height? :string ,
	width? :string ,
}



const MyButton = forwardRef<HTMLButtonElement,ButtonProps>(({height,width="100%",hoverBorder,hoverColor,activeBgColor="#1F22250F", activeColor ,isActive,disabled=false ,children ,role ,hoverBgColor,variant,colorVariant,onClickHandler, color,bgColor,text,padding,borderRadius,...restProps}:ButtonProps,ref)=>{

    return (
        <Button variant={variant} 
		ref={ref}
        role={role} 
        colorVariant={colorVariant} 
        color={color}
        bgColor={bgColor} 
        hoverBgColor={hoverBgColor}
        hoverColor={hoverColor}
        activeBgColor={activeBgColor}
        activeColor={activeColor}
        hoverBorder={hoverBorder}
        onClick={onClickHandler}
        padding={padding}
        borderRadius={borderRadius}
		height={height}
		isActive={isActive}
		width={width}
		disabled={disabled}

        {...restProps}
        >
            {children}
            {text}
        </Button>  
    )

})

export default MyButton

const Button = styled.button<{
    variant? : string ,
    role? : string ,
    colorVariant? : string ,
    color? : string ,
	activeBgColor? :string ,
	activeColor? :string ,
    bgColor? : string ,
    padding? : string  ,
    borderRadius? : string ; 
    height?:string ;
    width?:string ;
    hoverBgColor?:string ; 
    hoverBorder?:string ;
    hoverColor?:string ;
	isActive? :boolean ;
	disabled? :boolean ;

}>(props=>({
   
    backgroundColor : props.isActive && props.activeBgColor ||  (props.colorVariant == 'dark' ? 'black' : props.colorVariant == 'light' ? 'lightgrey' : props.colorVariant == 'medium'? 'grey' : props.role == 'danger'? 'red':  props.role == 'success'? 'green' :  props.role == 'primary'? 'blue': props.role =='primaryLight' ? props.theme.colors.button.primaryLight : props.bgColor ?? "inherit" ) ,
    color : props.color ?? "inherit" ,
    padding : props.padding ,
    borderRadius : props.borderRadius ?? "8px" ,
    border : 'none' ,
    boxShadow : 'none' ,
	outline : 'none' ,

	cursor: props.disabled ? 'not-allowed':'',
	inset : props.disabled ? 0 : undefined,
	opacity : props.disabled ? 0.65 : 1,

    boxSizing : 'border-box' ,
    '&:hover' :{
        boxSizing : 'border-box' ,
        backgroundColor :  props.isActive && props.activeBgColor || props.hoverBgColor ,
        color : props.hoverColor ,
        outline : props.hoverBorder ,
    },


}))
