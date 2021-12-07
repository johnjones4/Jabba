import { Component } from "react";
import JabbaEvent from "../../lib/Event";
import VendorType from "../../lib/VendorType";
import './events.css'

interface EventsProps {
  eventVendorType: string | null
  page: number
}

interface EventsState {
  events: Array<JabbaEvent>
  vendorTypes: Array<string>,
  page: string
}

export default class Events extends Component<EventsProps, EventsState> {
  constructor(props: EventsProps) {
    super(props)
    this.state = {
      events: [],
      vendorTypes: [],
      page: `${props.page}`
    }
  }

  async loadInfo() {
    try {
      this.setState({
        events: await JabbaEvent.loadEvents(this.props.eventVendorType, this.props.page),
        vendorTypes: await VendorType.loadVendorTypes(),
      })
    } catch (e) {
      console.error(e)
    }
  }
  
  componentDidMount() {
    this.loadInfo()
  }

  render() {
    return (
      <div className='JabbaEvents'>
        <form className='JabbaEvents-form'>
          <div className='JabbaEvents-form-item'>
            <select name='eventVendorType' onChange={e => e.target.form?.submit()}>
              { this.state.vendorTypes.map(t=> {
                return (
                  <option key={t} value={t} selected={t === this.props.eventVendorType}>{t}</option>
                )
              }) }
            </select>
          </div>
          <div className='JabbaEvents-form-item'>
            <input name='page' value={this.state.page} onChange={e => this.setState({page: e.target.value})} />
          </div>
        </form>
        <table>
          <thead>
            <tr>
              <th></th>
              <th>Date</th>
              <th>Type</th>
              <th>ID</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            { this.state.events.map((event, i) => {
              return (
                <tr key={i} className={event.isNormal ? 'normal' : 'abnormal'}>
                  <td>
                    <a href={`#/event/${event.id}`}>
                      View
                    </a>
                  </td>
                  <td>{event.created.toLocaleString()}</td>
                  <td>{event.eventVendorType}</td>
                  <td>{event.eventVendorID}</td>
                  <td>{event.isNormal ? 'Normal' : 'Abnormal'}</td>
                </tr>
              )
            }) }
          </tbody>
        </table>
      </div>
    )
  }
}
