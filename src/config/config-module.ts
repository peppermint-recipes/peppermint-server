import { Module } from '@nestjs/common';

import { ConfigService } from './config-service';

@Module({
    providers: [{
        provide: ConfigService,
        useFactory: () => new ConfigService(),
    }],
    exports: [ConfigService],
})
export class ConfigModule {}
