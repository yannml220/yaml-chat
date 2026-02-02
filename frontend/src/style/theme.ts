import { type Theme } from "@emotion/react";



export const appTheme : Theme = {
    colors : {
        background : { 
            darker : 'rgb(22 26 58 / 86%)',
            dark : '#495057' ,
            ligther : '#F8F9FA',
            ligth :'#d4d4de' ,
            neutral : '#000a200d',
            superlightgrey :'#1F22250F',
        },
        font : {
            dark : '#343A40' ,
            light : '#DEE2E6' ,
            indigo : '#2d4665' ,
            neutral : '#172b4d' ,
            soft : '#5e6c84' ,
        },
        button : {
            primaryLight :'#5199d7',
            primary : '#4361ee',
            secondary : '#00afb9',
            danger : '#ff0054',
            success : '#76c893',
        }
    },

    typography: {
        fontFamily: '-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Noto Sans,Ubuntu,Droid Sans,Helvetica Neue,sans-serif;',
        fontSizes:{
            xsmall : '0.8rem',
            small : '0.1rem',
            medium : '1.5rem',
            large : '2.5rem',
        },
        fontWeight: {
            light: 300,
            regular: 400,
            bold: 700,
        },
    }
    ,

    breakpoints: {
        sm: '640px',
        md: '768px',    
        lg: '1024px',
        xl: '1280px',
        xxl: '1536px'
      },

}


