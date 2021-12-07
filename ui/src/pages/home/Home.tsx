import React, { Component } from 'react';
import VendorType from '../../lib/VendorType';
import './home.css'

interface HomeProps {

}

interface HomeState {
  info: Array<VendorType>
}

export default class Home extends Component<HomeProps, HomeState> {
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
      <div className='Home'>
        { this.state.info.map((_info, i) => {
          return (
            <a href={`/events?eventVendorType=${_info.eventVendorType}`} key={i} className={`Home-status state-${_info.status}`}>
              <div className='Home-status-inner'>{ _info.eventVendorName }</div>
            </a>
          )
        }) }
      </div>
    )
  }

}
