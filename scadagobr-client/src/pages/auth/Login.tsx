import {Col, Row} from "antd";
import LoginForm from "../../components/auth/LoginForm";
import styled, {css} from "styled-components";
import {useUserStore} from "../../core/stores/userStore";
import {Navigate} from "react-location";
import React from "react";

const CenterForm = styled.div`
  margin-top: 30%;
  padding-left: 5%;
  padding-right: 5%;
`

const Title = styled.h1`
  font-size: 5rem;
  margin-bottom: 0;
  color: #F26419;
  font-family: "3270Medium Nerd Font Mono", serif;
  text-align: center;
`

function App() {
    const isLoggedIn = useUserStore(e => e.user.isLoggedIn)

    if (isLoggedIn) {
        return <Navigate to="/" />
    }

    return (
        <Row>
            <Col lg={14} sm={0} style={{backgroundColor: "#284b63"}}></Col>
            <Col lg={10} xs={24} style={{backgroundColor: "#FFFBFE", height: "100vh"}}>
                <CenterForm>
                    <Title>ScadaGOBR</Title>
                    <LoginForm></LoginForm>
                </CenterForm>
            </Col>
        </Row>
    )
}

export default App
