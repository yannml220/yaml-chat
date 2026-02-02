import React from 'react'
import styled from '@emotion/styled'

interface PageContentProps extends React.HTMLAttributes<HTMLDivElement> {
    children : React.ReactNode ,
    bgColor? : string | "white" ,
    justify? : string | undefined ,
    align? : string | undefined ,
    display? : string | 'block',
    padding? : string ;
	direction?:string ;
  }

const  PageContent = ({children,bgColor, justify , align ,display,padding,direction,...RestProps} : PageContentProps )=> {
  return (
    <Container 
    bgcolor ={bgColor} 
    justify={justify} 
    align={align} 
    display={display}
    padding={padding}
	direction={direction}
	{...RestProps}
    >
        {children}
    </Container>
  )
}

export default PageContent



const Container = styled.main<{
    bgcolor?:string ,
    justify?: string ,
    align? : string ,
    display? : string ,
    padding? : string  ,
	direction?:string ;
}>`
    min-height : 100% ;
    overflow-x : hidden  ; 
    box-sizing : border-box ;

    width : 100% ;
    display : ${(props)=>props.display} ;
    align-items : ${(props)=>props.align} ;
    justify-content : ${(props)=>props.justify} ;
    flex-direction : ${(props)=>props.direction} ;
    background-color : ${(props)=>props.bgcolor} ;
    padding : ${props=>props.padding};
`
