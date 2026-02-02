

interface ImgProps  extends React.HTMLAttributes<HTMLImageElement> {
	src : string,
	height : string,
	width : string ,
	alt? : string ,

}

const Img = ({src,height, width,alt,...Restprops}:ImgProps)=>{

	return (
		<img  src={src} height={height} width={width} alt={alt} {...Restprops}/>
	)


}


export default Img ;
