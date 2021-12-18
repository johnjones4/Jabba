import React, { Component } from 'react';
import { useSearchParams } from "react-router-dom"
import VendorType from '../../lib/VendorType';
import './home.css'
import './home-default.css'
import './home-lcars.css'
import { eventInfoString } from '../../lib/utils';

interface HomeProps {
  theme: string | null
}

interface HomeState {
  info: Array<VendorType>
}

class Home extends Component<HomeProps, HomeState> {
  private scroller: HTMLDivElement | null = null
  private scrollerInterval?: NodeJS.Timeout

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

  renderStatuses() {
    return (
      <div className='Home-statuses'>
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

  registerScroller() {
    if (!this.scrollerInterval) {
      const scrollerInterval = setInterval(() => {
        if (this.scroller) {
          let top = this.scroller.scrollTop + 10
          if (top > this.scroller.scrollHeight - this.scroller.clientHeight) {
            top = 0
            clearTimeout(scrollerInterval)
            this.scrollerInterval = undefined
            setTimeout(() => this.registerScroller(), 5000)
          }
          this.scroller.scrollTo({
            top: top,
            behavior: 'smooth'
          })
        }
      }, 100)
      this.scrollerInterval = scrollerInterval
    }
  }

  renderAbnormalStatusDetails() {
    const badStatuses = this.state.info.filter(info => info.status === 'abnormal')
    this.registerScroller()
    return (
      <div 
        className={['Home-bad-statuses-details', `Home-bad-statuses-details-${badStatuses.length > 0 ? 'nonempty' : 'empty'}`].join(' ')}
        ref={el => this.scroller = el}
      >
        { badStatuses.map((info, i) => {
          return (
            <div className='Home-bad-status-detail'>
              <h2>{ [info.eventVendorName, info.lastEvent.created.toLocaleString()].join(': ') }</h2>
              <pre>{ eventInfoString(info.lastEvent) }</pre>
            </div>
          )
        }) }
      </div>
    )
  }

  render() {
    return (
      <div className={['Home', `Home-theme-${this.props.theme ? this.props.theme : 'default'}`].join(' ')}>
        { this.renderStatuses() }
        { this.renderAbnormalStatusDetails() }
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
