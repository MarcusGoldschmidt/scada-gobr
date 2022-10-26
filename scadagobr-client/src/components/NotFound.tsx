import NotFoundLogo from "../assets/undraw_just_browsing_re_ofnd.svg";
import { Typography } from 'antd';

export default function NotFound() {
    return (
        <div style={{textAlign: "center", marginTop: "20vh"}}>
            <Typography.Title>Page not found 😥</Typography.Title>
            <img src={NotFoundLogo} alt="Not found image"/>
        </div>
    )
}