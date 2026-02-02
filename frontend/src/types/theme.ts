import '@emotion/react'; // it's important to use ThemeProvider
declare module '@emotion/react' {
  export interface Theme {
    colors : {
        background : { 
            darker : string,
            dark : string ,
            ligther : string,
            ligth :string ,
            neutral : string ,
            superlightgrey : string
            
        },
        font : {
            dark : string ,
            light : string ,
            indigo:string ,
            neutral : string ,
            soft : string ,
        },
        button : {
            primaryLight : string ,
            primary : string,
            secondary : string,
            danger : string,
            success : string,
        }
    },

    typography: {
        fontFamily: string,
        fontSizes:{
            xsmall : string,
            small : string,
            medium : string,
            large : string,
        },
        fontWeight: {
            light: number,
            regular: number,
            bold: number,
        },
    }
    ,

    breakpoints: {
        sm: string,
        md: string,    
        lg: string,
        xl: string,
        xxl: string
      },
  }
}
