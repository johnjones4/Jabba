import { Component } from "react";
import JabbaEvent from "../../lib/Event";
import { eventInfoString } from "../../lib/utils";
import './event.css'

interface EventProps {
  id: number
}

interface EventState {
  event: JabbaEvent | null
}

export default class Event extends Component<EventProps, EventState> {
  constructor(props: EventProps) {
    super(props)
    this.state = {
      event: null
    }
  }

  async loadInfo() {
    try {
      this.setState({
        event: await JabbaEvent.load(this.props.id)
      })
    } catch (e) {
      console.error(e)
    }
  }
  
  componentDidMount() {
    this.loadInfo()
  }

  render() {
    if (!this.state.event) {
      return null
    }
    const data = [
      ['ID', this.state.event.id],
      ['Type', this.state.event.eventVendorType],
      ['Vendor ID', this.state.event.eventVendorID],
      ['Created', this.state.event.created.toLocaleString],
      ['Status', this.state.event.isNormal ? 'Normal' : 'Abnormal'],
      ['Info', (<pre>{eventInfoString(this.state.event)}</pre>)],
    ]
    return (
      <div className='JabbaEvent'>
        <table>
          <thead>
            <tr>
              <th>Label</th>
              <th>Value</th>
            </tr>
          </thead>
          <tbody>
            { data.map(([k, v]) => {
              return (
                <tr>
                  <td className='key'>{k}</td>
                  <td className='value'>
                    <pre>{v}</pre>
                  </td>
                </tr>
              )
            }) }
          </tbody>
        </table>
      </div>
    )
  }
}
