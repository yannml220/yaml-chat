
interface AvatarProps  extends React.HTMLAttributes<HTMLDivElement>  {
	children : React.ReactNode ;

}


const Avatar = ({children,...Restprops}: AvatarProps) =>{

	return (
		<div
		style={{
			display:"flex",
			justifyContent:"center",
			alignItems:"center",
			height : "35px",
			width : "35px",
			borderRadius : "50%" ,

		}}
		{...Restprops}

		>
		 {children}
		</div>


	)


}



export default Avatar
