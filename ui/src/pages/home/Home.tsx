import React, { Component } from 'react';
import { useSearchParams } from "react-router-dom"
import VendorType from '../../lib/VendorType';
import './home.css'
import './home-default.css'
import './home-lcars.css'

interface HomeProps {
  theme: string | null
}

interface HomeState {
  info: Array<VendorType>
}

class Home extends Component<HomeProps, HomeState> {
  constructor(props: HomeProps) {
    super(props)
    this.state = {
      info: []
    }
  }

  async reload() {
    try {
      this.setState({
        info: await VendorType.loadAll()
      })
    } catch (e) {
      console.error(e)
    }
  }

  async componentDidMount() {
    this.reload()
    setInterval(async () => {
      this.reload()
    }, 60000)
  }

  render() {
    return (
      <div className={['Home', `Home-theme-${this.props.theme ? this.props.theme : 'default'}`].join(' ')}>
        { this.state.info.map((_info, i) => {
          return (
            <a href={`#/events?eventVendorType=${_info.eventVendorType}`} key={i} className={`Home-status state-${_info.status}`}>
              <div className='Home-status-inner'>{ _info.eventVendorName }</div>
            </a>
          )
        }) }
      </div>
    )
  }

}

const HomeWrapper = () => {
  const [params] = useSearchParams()
  const theme = params.get('theme')
  return (<Home theme={theme} />)
}

export default HomeWrapper
